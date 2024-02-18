# Broker Mail (Send Grid)

BrokerMail is a simple HTTP server written in Go that facilitates sending emails using the SendGrid API. The project aims to focus on understanding communication logic and underlying concepts rather than adhering strictly to Go's architectural principles or leveraging SendGrid's utility packages.

## Developer Notes
As the developer of BrokerMail, I'd like to highlight that the primary goal of this project is not to showcase advanced architectural patterns in Go or extensively utilize SendGrid's utility packages. Instead, the emphasis is on grasping fundamental communication logic and concepts under the hood.

While the code may not follow all best practices or utilize all available tools provided by SendGrid, it serves as an educational tool to understand the basics of handling HTTP requests, structuring data, and interacting with external APIs.

## Prerequisites
Before running this application, ensure you have the following:

- SendGrid API Key
- SendGrid API URL
- Sender email address
- Port number to run the server on

## Setup
1. Clone this repository.
2. Create a **.env** file in the root directory of the project.
3. Add the following environment variables to the **.env** file:
```.env
SENDGRID_API_KEY=your_sendgrid_api_key
SENDGRID_API_URL=sendgrid_api_url
SENDER=your_sender_email_address
PORT=port_number
```
4. Install required dependencies using go mod tidy.
5. Run the application using go run main.go.

## Usage
To send an email, make a POST request to the **/send** endpoint with the following JSON payload:

```json
{
  "subject": "Your Email Subject",
  "content": "Your Email Content",
  "recipients": ["recipient1@example.com", "recipient2@example.com"]
}
```
Replace subject, content, and recipients with your desired email details.

## Architecture
The main components of the application include:

- **main.go**: The entry point of the application. It initializes the server and handles incoming requests.
- **SendGridConfig**: Struct holding SendGrid configuration and implementing the ServeHTTP method to handle requests.
- **models**: Package containing data structures for request and response bodies.
- **utils**: Package containing utility functions for validating JSON payloads.
  
## Contributing
Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

