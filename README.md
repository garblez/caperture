# Portfolio Website for Jonathan
Simple portfolio website for Jonathan that requires no dependencies and runs on a Go+HTMX stack.

Deployment to the server will be done via git pushes to a yet-to-be-defined remote repo on the server itself with a custom deployment script run after a merge hook is triggered.
NGINX should then load the binary as a systemd daemon and run the html in-memory allowing for efficient caching.

### Environment keys
This just uses a simple net/smtp call to send the email to a specified email address. This should allow hiding private emails behind a customer contact form. The following need to be defined as environment variables when running the server:
* `GSMTP_EMAIL` The email address with which the email shall be sent.
* `GSMTP_RECIPIENT` The email address that should receive the email.
* `GSMTP_PASSWORD` The app-specific password for the `GSMTP_EMAIL` account

If any of these are undefined, an error will occur on the backend and will be logged (as of now, there are no persistent log files and errors are just written to the standard Logger output.) The customer-facing frontend currently just displays the success response fragment (thanks.html) irrespective of outcome.

### Deployment
Deployment is done by git pushing to the `deploy` remote repository. The post-receive hook should then fire off a bash script for building the project in /srv/tmp and moving it (on success) /src/www.

NGINX is used as the server of choice to provide the application.
