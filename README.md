# ketu_backend_monolith_v1

## 📊 Database Configuration

### Connection Pool Settings

- **maxOpenConns:** 25 (default)

  - Formula: (CPU cores \* 2) + effective_spindle_count
  - Monitor connection wait times to adjust

- **connMaxLifetime:** 15 minutes

  - Keep lower than database server timeout
  - Adjust if using connection proxy

- **maxIdleConns:** 25

  - Equal to maxOpenConns for busy applications
  - Can be lower for less busy systems

- **connMaxIdleTime:** 10 minutes
  - Adjust based on traffic patterns
  - Lower if memory is constrained

## 🔒 Security

- JWT-based authentication
- Password hashing with bcrypt
- Role-based access control
- Request validation
- SQL injection protection

## 📝 API Documentation

Full API documentation is available at `/swagger/index.html` when running the application.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
