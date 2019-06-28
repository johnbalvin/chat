class SnackBar{
  constructor(){
    this.me="";
    this.img="";
    this.texto="";
    this.simbCargando;
    this.simbCross;
    this.simbCheck;
    this.mostrado=false;
    this.worker="";
    this.genWorker();
    this.mostrar=this.mostrar.bind(this);
    this.desaparecer=this.desaparecer.bind(this);
    this.desaparecer2=this.desaparecer2.bind(this);
    this.cambiarSimbolo=this.cambiarSimbolo.bind(this);
    this.cambiarMens=this.cambiarMens.bind(this);
    this.cambiarMensaYSimb=this.cambiarMensaYSimb.bind(this);
    this.menSimbTime=this.menSimbTime.bind(this);
    this.start();
  }
  start(){
    const me=`
    <div id="snackbar">
      <h2 class="snackbartext">Cargando....</h2>
      <div class="snackbarImg"></div>
    </div>
    `;
    const simbCargando=`
    <svg class="circular-loader" width="50" height="50" viewBox="25 25 50 50" >
      <circle class="loader-path" cx="50" cy="50" r="20" fill="none" stroke="#70c542" stroke-width="10" />
    </svg>
    `;
    const simbCross=`
    <svg width="50" height="50" viewBox="0 0 32 32">
     <path d="M31.708 25.708c-0-0-0-0-0-0l-9.708-9.708 9.708-9.708c0-0 0-0 0-0 0.105-0.105 0.18-0.227 0.229-0.357 0.133-0.356 0.057-0.771-0.229-1.057l-4.586-4.586c-0.286-0.286-0.702-0.361-1.057-0.229-0.13 0.048-0.252 0.124-0.357 0.228 0 0-0 0-0 0l-9.708 9.708-9.708-9.708c-0-0-0-0-0-0-0.105-0.104-0.227-0.18-0.357-0.228-0.356-0.133-0.771-0.057-1.057 0.229l-4.586 4.586c-0.286 0.286-0.361 0.702-0.229 1.057 0.049 0.13 0.124 0.252 0.229 0.357 0 0 0 0 0 0l9.708 9.708-9.708 9.708c-0 0-0 0-0 0-0.104 0.105-0.18 0.227-0.229 0.357-0.133 0.355-0.057 0.771 0.229 1.057l4.586 4.586c0.286 0.286 0.702 0.361 1.057 0.229 0.13-0.049 0.252-0.124 0.357-0.229 0-0 0-0 0-0l9.708-9.708 9.708 9.708c0 0 0 0 0 0 0.105 0.105 0.227 0.18 0.357 0.229 0.356 0.133 0.771 0.057 1.057-0.229l4.586-4.586c0.286-0.286 0.362-0.702 0.229-1.057-0.049-0.13-0.124-0.252-0.229-0.357z" fill="red"></path>
    </svg>
    `;
    const simbCheck=`
    <svg width="50" height="50" viewBox="0 0 32 32">
      <path d="M27 4l-15 15-7-7-5 5 12 12 20-20z" fill="green"></path>
    </svg>
    `
    let template=document.createElement("template");
    template.innerHTML=me;
    this.me=template.content.querySelector("#snackbar");
    document.body.appendChild(this.me);
    template.innerHTML=simbCargando;
    this.simbCargando=template.content.querySelector("svg");
    template.innerHTML=simbCross;
    this.simbCross=template.content.querySelector("svg");
    template.innerHTML=simbCheck;
    this.simbCheck=template.content.querySelector("svg");
    this.img=this.me.querySelector(".snackbarImg");
    this.img.appendChild(this.simbCargando);
    this.texto=this.me.querySelector(".snackbartext");
  }
  mostrar(mensaje,simbolo,tiempo="nop"){
    if(this.mostrado){
      this.worker.addEventListener("message",()=>{this.mostrar(mensaje,simbolo,tiempo)},{once:true});
      return
    }
    this.mostrado=true;
    this.img.firstChild.replaceWith(simbolo);
    this.texto.textContent=mensaje;
    this.me.style.visibility="visible";
    if(tiempo!="nop"){this.desaparecer2(tiempo)}
    this.me.style.transform="translate(-50%,-50%)";
  }
  genWorker(){
    const content = "onmessage=function(mess){postMessage(mess.data[0]);}";
    const blob = new Blob([content],{type:'application/javascript'});
    this.worker=new Worker(URL.createObjectURL(blob));
  }
  desaparecer(tiempo){
    setTimeout(()=>{
      this.me.style.transform="translate(-50%,55vh)";
      this.me.addEventListener("transitionend",()=>{
        this.me.style.visibility="hidden";
        this.mostrado=false;
        this.worker.postMessage(["ended"])},{once:true});
    },tiempo*1000)
  }
  desaparecer2(tiempo){
    this.me.addEventListener("transitionend",()=>{
        setTimeout(()=>{
          this.me.style.transform="translate(-50%,55vh)";
          this.me.addEventListener("transitionend",()=>{
            this.me.style.visibility="hidden";
            this.mostrado=false;
            this.worker.postMessage(["ended"]);
        },{once:true});
        },tiempo*1000) 
    },{once:true})
  }
  cambiarSimbolo(simbolo){
    this.img.firstChild.replaceWith(simbolo);
  }
  cambiarMens(mensaje){
    this.texto.textContent=mensaje;
  }
  cambiarMensaYSimb(mensaje,simbolo){
    this.texto.textContent=mensaje;
    this.img.firstChild.replaceWith(simbolo);
  }
  menSimbTime(mensaje,simbolo,tiempo){
    this.texto.textContent=mensaje;
    this.img.firstChild.replaceWith(simbolo);
    this.desaparecer(tiempo);
  }
}
let snackBar=new SnackBar();
