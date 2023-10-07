# noveau_galaxy
This is the Core Banking Application Layer of Nouveau Mobile

# **Account Creation API**

The Account Creation API is a crucial component of a banking system that allows users to open new accounts. It collects user information and creates a new account, associating it with a unique account number. Below is an explanation of the request, response, and the logic behind this API.

## **Request**

The request represents the data sent by a user or a client application when they want to open a new account. It includes various pieces of information needed for the account creation process.

- **`surname`**: The user's last name or surname.
- **`firstName`**: The user's first name.
- **`homeAddress`**: The user's home address.
- **`city`**: The city where the user resides.
- **`phoneNumber`**: The user's contact phone number.
- **`identity`**: An object containing identity-related information.
    - **`bvn`**: The user's Bank Verification Number.
    - **`passport`**: The user's passport information.
    - ``ut**ilityBill**`: A utility bill as proof of address.
- **`picture`**: A required field for the user's picture or photograph.
- **`country`**: The user's country of residence.
- date : The user's country of residence.
- identity : id as foreign key

## **Response**

The response represents the result of the account creation process. It includes various pieces of information that provide feedback to the user or client application.

- **`responseCode`**: A code indicating the result of the operation (e.g., "00" for success).
- **`responseMessage`**: A message describing the outcome of the operation (e.g., "Successful").
- **`accountNumber`**: The unique account number assigned to the newly created account.
- **`custNumber`**: A customer number or identifier associated with the user.
- **`imagePt`**: An identifier for the user's picture or image.

## **Logic**

The logic behind the Account Creation API involves several steps:

1. **Data Collection**: The API collects user-provided information, including personal details, contact information, identity documents (e.g., BVN, passport, utility bill), and a  picture.
2. **Storage in Users Table**: The collected user data is stored in the Users table, which acts as a repository for user information. This information is associated with a unique identifier.
3. **Account Number Generation**: A new account number is generated for the user. This number is typically unique and serves as an identifier for the user's account within the bank.
4. **Storage in Accounts Table**: The generated account number is stored in the Accounts table, and it is associated with a user identifier (user-id). It's important to note that a single user can have multiple account numbers associated with their user ID, allowing them to manage multiple accounts.
5. **Storage in Identification Table:** The notification json request data is going to be stored in a notification table and associated with a user-id / account  number 

## Databases

- - Transactions Table
CREATE TABLE Transactions (
transaction_id INT PRIMARY KEY,
sender_account_id INT,
receiver_account_id INT,
amount DECIMAL(10,2),
currency_code VARCHAR(3),
status ENUM('pending', 'completed', 'failed'),
transaction_type ENUM('deposit', 'withdrawal', 'transfer'),
timestamp TIMESTAMP,
FOREIGN KEY (sender_account_id) REFERENCES Accounts(account_id),
FOREIGN KEY (receiver_account_id) REFERENCES Accounts(account_id)
);
- - Beneficiaries Table
CREATE TABLE Beneficiaries (
beneficiary_id INT PRIMARY KEY,
user_id INT,
full_name VARCHAR(100),
bank_name VARCHAR(100),
bank_account_number VARCHAR(20),
bank_routing_number VARCHAR(20),
swift_code VARCHAR(20),
created_at TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES Users(user_id)
);
- - KYC Documents Table
CREATE TABLE KYC_Documents (
document_id INT PRIMARY KEY,
user_id INT,
document_type ENUM('passport', 'national ID'),
document_number VARCHAR(20),
document_image_path VARCHAR(255),
expiry_date DATE,
FOREIGN KEY (user_id) REFERENCES Users(user_id)
);
- - Compliance Logs Table
CREATE TABLE Compliance_Logs (
log_id INT PRIMARY KEY,
user_id INT,
action VARCHAR(100),
timestamp TIMESTAMP,
details TEXT,
FOREIGN KEY (user_id) REFERENCES Users(user_id)
);


## Account Generation Service
**`Background`**:In the context of banking in the United Kingdom (UK), an account number and a sort code are two crucial pieces of information associated with a bank account.

- Account Number:

The account number is a unique identifier assigned to an individual or entity's bank account. It is used to distinguish one account from another within a specific bank or financial institution.
In the UK, traditional account numbers typically consist of six digits. These numbers are usually assigned sequentially, starting from a specific range.
The account number helps in routing transactions to the correct account within the bank.

- Sort Code:

The sort code is another essential element in the UK banking system. It identifies the specific branch or office of a bank where an account is held.
Sort codes consist of six digits as well.
The sort code is used in conjunction with the account number to ensure that a transaction is directed to the correct branch and account.
Together, the account number and sort code work in tandem to facilitate various financial transactions, including deposits, withdrawals, and transfers. When combined, they provide a unique address for each bank account within the UK banking system. This combination is used by the banking infrastructure to accurately route payments and transfers to the intended recipient.

- **About our directory structure**

```
- ukaccountgen
      - ukaccountgen.go
- main.go(Execution or test file)
```
