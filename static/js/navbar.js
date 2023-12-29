//console.log("404", "NOT FOUND!");

function myFunction() {
    console.log("RESIZE!!!");
    var x = document.getElementById("myTopnav");
    if (x.className === "topnav") {
        x.className += " responsive";
    } else {
        x.className = "topnav";
    }
}