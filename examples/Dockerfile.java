# Java application with multistage build but other issues
FROM openjdk:17 AS builder

WORKDIR /build

# Copy pom first for better caching
COPY pom.xml .
# Use wildcard copy for src
COPY src ./src

# Build with Maven
RUN ./mvnw package -DskipTests

# Second stage but still using large base image
FROM openjdk:17

WORKDIR /app

# Copy JAR from builder stage
COPY --from=builder /build/target/*.jar app.jar

# Expose port
EXPOSE 8080

# Set healthcheck
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/actuator/health || exit 1

# Still running as root (no USER)

# Run application
CMD ["java", "-jar", "app.jar"]