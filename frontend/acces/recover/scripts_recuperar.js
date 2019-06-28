class Recuperar{
  constructor(){
    this.handlearSubmit=this.handlearSubmit.bind(this);
    this.me=document.querySelector("#recuperar");
    this.paso1=document.querySelector("#paso1");
    this.correoInput=this.paso1.querySelector("input");
    this.paso2=document.querySelector("#paso2");
    this.codigoInput=this.paso2.querySelector("input");
    this.paso3=document.querySelector("#paso3");
    this.contraseñaInput=this.paso3.querySelector(".contra");
    this.paso=1;
    this.comenzar();
  }
  comenzar(){
    this.me.addEventListener("submit",this.handlearSubmit);
  }
  handlearSubmit(e){
    e.preventDefault();
    let correo=this.correoInput.value;
    let codigo=this.codigoInput.value;
    let contraseña=this.contraseñaInput.value;
    let send=new FormData();
    send.append("correo",correo);
    send.append("p",this.paso);
    snackBar.mostrar("Cargando...",snackBar.simbCargando);
    switch (this.paso){
      case 1:
        fetch("",{method:"POST",body:send,credentials:'include'})
        .then(resp =>{
          switch (resp.status){
            case 409:
              snackBar.menSimbTime("El usuario ya esta registrado",snackBar.simbCross,5);
              break;
            case 500:
              snackBar.menSimbTime("Error del servidor",snackBar.simbCross,5);
              break;   
            case 200:
              this.paso1.style.display="none";
              this.paso2.style.display="block";
              this.paso3.style.display="none";
              this.codigoInput.focus();
              this.paso=2;
              snackBar.menSimbTime("Continue",snackBar.simbCheck,1);
              break;  
          }
        }).catch(err=> {
          console.log(err);
          snackBar.menSimbTime("Error de conexión, intente de nuevo",snackBar.simbCross,3)
        });
        break;
      case 2:
        send.append("code",codigo);
        fetch("",{method:"POST",body:send,credentials:'include'})
        .then(resp =>{
          switch (resp.status){
            case 409:
              snackBar.menSimbTime("El código expiró".snackBar.simbCross,5);
              break;
            case 500:
              snackBar.menSimbTime("Error del servidor",snackBar.simbCross,5);
              break;   
            case 200:
              this.paso1.style.display="none";
              this.paso2.style.display="none";
              this.paso3.style.display="block";
              this.contraseñaInput.focus();
              this.paso=3;
              snackBar.menSimbTime("Continue",snackBar.simbCheck,1);
              break;  
          }
        }).catch(err=> {
          console.log(err);
          snackBar.menSimbTime("Error de conexión, intente de nuevo",snackBar.simbCross,4)
        });
        break;
      case 3:
        send.append("code",codigo);
        send.append("pass",contraseña);
        fetch("",{method:"POST",body:send,credentials:'include'})
        .then(resp =>{
          switch (resp.status){
            case 409:
              snackBar.menSimbTime("El código expiró".snackBar.simbCross,5);
              break;
            case 500:
              snackBar.menSimbTime("Error del servidor",snackBar.simbCross,5);
              break;   
            case 200:
              window.location.href = "/MisSuperCursos"; 
              break;
          }
        }).catch(err=> {
          console.log(err);
          snackBar.menSimbTime("Error de conexión, intente de nuevo",snackBar.simbCross,2)
        });
        break;    
    }
  }
}
let recuperar=new Recuperar();