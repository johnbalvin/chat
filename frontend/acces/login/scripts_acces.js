class Form{
  constructor(){
    this.me=document.querySelector("#myform");
    this.correoInput=document.querySelector("#email");
    this.contraInput=document.querySelector("#pass");
    this.sendToserver=this.sendToserver.bind(this);
    this.start();
  }
  start(){
    this.me.addEventListener("submit",this.sendToserver);
  }
  sendToserver(e){
    e.stopPropagation();e.preventDefault();
    let send=new FormData();
    send.append("email",this.correoInput.value);
    send.append("pass",this.contraInput.value);
    snackBar.mostrar("Waiting...",snackBar.simbCargando);
    fetch("",{method:"POST",body:send,credentials:'include'})
    .then(resp =>{
      switch (resp.status){
        case 404:
          snackBar.menSimbTime("User not in database",snackBar.simbCross,5);
          break;
        case 500:
          snackBar.menSimbTime("Server error",snackBar.simbCross,5);
          break;
        case 422: 
          snackBar.menSimbTime("Wrong password",snackBar.simbCross,5);
          break;
        case 200:
          location.reload();   
          break;  
      }
    }).catch(err=> {
      console.log(err);
      snackBar.menSimbTime("Connection error, try again",snackBar.simbCross,2)
    });
  }
}
new Form(); 