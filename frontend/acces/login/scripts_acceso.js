class Form{
  constructor(){
    this.hide=this.hide.bind(this);
    this.me=document.querySelector("#login");
    this.dark=document.querySelector("#dark");
    this.correoInput=document.querySelector("#correo");
    this.contraInput=document.querySelector("#contra");
    this.enviarServidor=this.enviarServidor.bind(this);
    this.comenzar();
  }
  comenzar(){
    this.me.addEventListener("submit",this.enviarServidor);
  }
  enviarServidor(e){
    e.stopPropagation();e.preventDefault();
    let send=new FormData();
    send.append("email",this.correoInput.value);
    send.append("pass",this.contraInput.value);
    snackBar.mostrar("Waiting...",snackBar.simbCargando);
    fetch("",{method:"POST",body:send,credentials:'include'})
    .then(resp =>{
      switch (resp.status){
        case 404:
          snackBar.menSimbTime("El usuario no se encuentra en la base de datos",snackBar.simbCross,5);
          break;
        case 500:
          snackBar.menSimbTime("Error del servidor",snackBar.simbCross,5);
          break;
        case 422: 
          snackBar.menSimbTime("Clave incorrecta",snackBar.simbCross,5);
          break;
        case 200:
          location.reload();   
          break;  
      }
    }).catch(err=> {
      console.log(err);
      snackBar.menSimbTime("Error de conexión, intente de nuevo",snackBar.simbCross,2)
    });
  }
}
new Form(); 