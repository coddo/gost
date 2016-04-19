package email

const activateAccountTemplate = `<html>
<body>
    <h1> Hello %s</h1> <br><br>
    <h2> Welcome to the GOST web framework!</h2>
    <br>
    <br>
    <p style="text-size: 20px">To activate your account, please use the following link: %s</p>
</body>
</html>`

const resetPasswordTemplate = `<html>
<body>
    <h1> Hello %s</h1> <br><br>
    <h2>A password reset was requested for this account. If it was not you who made the request, please disregard this email</h2>
    <br>
    <br>
    <p style="text-size: 20px">To reset your password, please use the following link: %s</p>
</body>
</html>`
