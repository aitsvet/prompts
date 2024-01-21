FROM gradle:jdk17 as build_env
WORKDIR /app
COPY ./build.gradle .
COPY . .
RUN gradle clean bootJar -x test

FROM openjdk:17-slim
WORKDIR /app
COPY --from=build_env /app/build/libs/*.jar app.jar
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-jar","/app/app.jar"]
