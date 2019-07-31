const express = require('express');
const app = express();
const router = express.Router();

const views = __dirname + '/views/';
const port = 1234;

router.use(function (req,res,next) {
    console.log('/' + req.method + " " + req.path);
    next();
});

router.get('/', function(req,res){
    res.sendFile(views + 'index.html');
});

app.use(express.static(views));
app.use('/', router);

app.listen(port, function () {
  console.log('App listening on port ' + port + "!")
})
