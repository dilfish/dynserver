package main

var MessagePage = `
<!doctype html>
<html lang="zh-cmn-Hans">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">

    <title>
input message
</title>

  </head>
  <body>

  <form action="/t" name="confirmationForm" method="post">
    <textarea id="message" class="text" cols="40" rows ="20" name="message"></textarea>

   <input type="submit" value="send" class="submitButton">
</form>


  </body>
</html>
`
