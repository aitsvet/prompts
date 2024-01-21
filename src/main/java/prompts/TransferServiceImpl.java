package prompts;

import balance.*;
import balance.BalanceOuterClass.*;
import transfer.*;
import transfer.TransferOuterClass.*;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.stub.StreamObserver;
import org.lognet.springboot.grpc.GRpcService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.core.env.Environment;

@GRpcService
public class TransferServiceImpl extends TransferGrpc.TransferImplBase {

    @Autowired
    private Environment env;

    private ManagedChannel channel1, channel2;

    public TransferServiceImpl(Environment env) {
        this.env = env;
        String url1 = env.getProperty("balance.0.url");
        channel1 = ManagedChannelBuilder.forTarget(url1).usePlaintext().build();
        String url2 = env.getProperty("balance.1.url");
        channel2 = ManagedChannelBuilder.forTarget(url2).usePlaintext().build();
    }

    @Override
    public void transferMoney(TransferRequest request, StreamObserver<TransferResponse> responseObserver) {
        int fromInstance = (int)request.getFromAccountId() % 2;
        int toInstance = (int)request.getToAccountId() % 2;
        
        BalanceGrpc.BalanceBlockingStub balanceClient1 = 
            BalanceGrpc.newBlockingStub(fromInstance == 0 ? channel1 : channel2);
        
        BalanceGrpc.BalanceBlockingStub balanceClient2 = 
            BalanceGrpc.newBlockingStub(toInstance == 0 ? channel1 : channel2);
        
        ChangeRequest changeFrom = ChangeRequest.newBuilder()
                 .setTransactionId(System.currentTimeMillis())
                 .setAccountId(request.getFromAccountId())
                 .setAmount(-request.getAmount())
                 .build();
        balanceClient1.changeBalance(changeFrom);
        
        ChangeRequest changeTo = ChangeRequest.newBuilder()
                .setTransactionId(System.currentTimeMillis())
                .setAccountId(request.getToAccountId())
                .setAmount(request.getAmount())
                .build();
        balanceClient2.changeBalance(changeTo);
        
        TransferResponse response = TransferResponse.newBuilder().setSuccess(true).build();
        responseObserver.onNext(response);
        responseObserver.onCompleted();
     }

    @Override
    public void getAccountData(TransferOuterClass.GetDataRequest request, StreamObserver<TransferOuterClass.GetDataResponse> responseObserver) {
        int instance = (int)request.getAccountId() % 2;
        
        BalanceGrpc.BalanceBlockingStub balanceClient = 
            BalanceGrpc.newBlockingStub(instance == 0 ? channel1 : channel2);
        
            BalanceOuterClass.GetDataRequest dataRequest = BalanceOuterClass.GetDataRequest.newBuilder()
                .setAccountId(request.getAccountId())
                .build();
                BalanceOuterClass.GetDataResponse balanceResponse = balanceClient.getAccountData(dataRequest);

                var b = TransferOuterClass.GetDataResponse.newBuilder();
                b.setBalance(balanceResponse.getBalance());
                for (var t : balanceResponse.getTransactionsList()) {
                    var tb = TransferOuterClass.GetDataResponse.Transaction.newBuilder();
                    tb.setId(t.getId());
                    tb.setAmount(t.getAmount());
                    b.addTransactions(tb.build());
                }
        
        responseObserver.onNext(b.build());
        responseObserver.onCompleted();
    }
}