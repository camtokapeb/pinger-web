//alert("НАЧАЛИ!!!!!!!!!!!!!!!!!!!!!!!!!");

const http = new XMLHttpRequest();

http.open("GET", "https://10.184.254.32:8765/books/www", true);

//http.onload = () => console.log(http.responseText)

//console.log(http);
http.onload = function(){
    if(http.status == 200){
        console.log("200", "OK!");
        console.log(http.responseText);
    } else if(http.status == 404){
        console.log("404", "NOT FOUND!");
    }
  };


  
  //resp = '{\"id\":2,\"title\":\"Преступление и наказание!!!!\"}';
  //console.log(JSON.parse(resp).title);
  
  http.send();

//alert("ЗАКОНЧИЛИ!");