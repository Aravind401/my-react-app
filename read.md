#   Please install the latest version for golang and react
# golang https://go.dev/doc/install
# react installation guide https://www.freecodecamp.org/news/how-to-install-react-a-step-by-step-guide/

#   1. Installing Dependencies
#   Install API dependencies:
#   Run the PowerShell script FE.ps1.

#   Install Front-End dependencies:
#   Run the PowerShell script BE.ps1.

#   This will create two folders: my-react-app (front-end) and myapi (API), each containing the necessary          dependencies.

#   2. Pre-requisites
#   Before running the application, ensure the following versions are installed:

#   Go Version:
#   Open Command Prompt and run:

#   cmd
#   go version
#   Expected output:
#   go version go1.24.4 windows/amd64


#   Node.js Version:
#   cmd
#   node --version
#   Expected output:
#   v18.16.0

#   npm Version:
#   cmd
#   npm --version
#   Expected output:
#   9.5.1

#   3. Running the Applications
#   Front-End (React App)
#   Open Command Prompt.
#   Navigate to the my-react-app directory.
#   Run the following command:
#   cmd 
#   npm start
#   open the url you can see the login page http://localhost:3000/


#   API (Go App)
#   Open Command Prompt.
#   Navigate to the myapi directory.
#   Run the following command:
#   cmd
#   go run main.go
#   open the url you can see the api documentation in http://localhost:8080/swagger/index.html

# setting up postgress sql server have to download windows installer from official the version we are using 17.5
# https://www.enterprisedb.com/downloads/postgres-postgresql-downloads


# once postgress sql is done just login to postgress sql and execute the below cmd to create the table following table:
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price NUMERIC(10,2) NOT NULL,
    in_stock BOOLEAN DEFAULT TRUE
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
