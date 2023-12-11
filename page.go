package main

var MessagePage = `
<!doctype html>
<html lang="zh-cmn-Hans">
  <head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <meta name="description" content="">
    <meta name="author" content="">
    <link rel="stylesheet" href="https://cdn.871116.xyz/302/bulma.css">


    <title>
        输入内容
</title>

  </head>
  <body>


<div class="columns">
  <div class="column is-full">
<article class="message">
  <div class="message-body">
      你可以输入任意内容，也可以从 curl 提交：curl -X POST -d 'message=测试从curl提交' https://dev.ug/t
  </div>
</article>
  </div>
</div>


<div class="columns is-mobile">
  <div class="column is-10 is-offset-1">

  <form action="/t" name="confirmationForm" method="post" align="center">
   <textarea class="textarea" placeholder="例如：我要记的网址是 https://dev.ug" name="message"></textarea>
   <button class="button is-primary" value="send">提交</button>
 </form>
 </div>


</div>


  </body>
</html>
`
