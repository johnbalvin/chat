const nameInput=document.querySelector(".configuracionNombre_input");
class UserImg{
  constructor(){
    this.me=document.querySelector(".configuracionFoto_img");
    this.inputBtn=document.querySelector(".configuracionFoto_boton");
    this.inputImg=document.querySelector(".configuracionFoto_input");
    this.src=this.me.src;
  }
} 
let userImg=new UserImg();
class Actions{
  constructor(){
    this.cancelBtn=document.querySelector(".cancelBotoncito");
    this.enviarBtn=document.querySelector(".enviarBotoncito");
    this.cancelar=this.cancelar.bind(this);
    this.enviar=this.enviar.bind(this);
    this.showAll=this.showAll.bind(this);
    this.hideAll=this.hideAll.bind(this);
    this.makePreview=this.makePreview.bind(this);
    this.name1=document.querySelector("#hereName");
    this.name2=document.querySelector("#nameMenu1");
    nameInput.addEventListener("input",this.showAll,{once:true})
    this.cancelBtn.addEventListener("click",this.cancelar);
    this.enviarBtn.addEventListener("click",this.enviar);
    userImg.inputImg.addEventListener("change",this.makePreview);
    userImg.inputBtn.addEventListener("click",()=>{userImg.inputImg.click()});
  }
  makePreview(e){
    botonsActn.showAll(e);
    userImg.me.src=window.URL.createObjectURL(userImg.inputImg.files[0]);
  }
  cancelar(e){
      this.hideAll(e);
      nameInput.addEventListener("input",this.nameInputListener,{once:true})
      userImg.me.src=userImg.src;
      nameInput.value=nameInput.dataset.value;
  }
  enviar(e){
    this.hideAll(e);
    let send=new FormData(),noEnviar=true;
    if(userImg.me.src!=userImg.src){send.append("i",userImg.inputImg.files[0]);noEnviar=false}
    if(nameInput.value!=nameInput.dataset.value){
      nameInput.dataset.value=nameInput.value;
      this.name1.textContent=nameInput.value;
      this.name2.textContent=nameInput.value;
      send.append("n",nameInput.value);noEnviar=false;
    }
    if(noEnviar)return
    snackBar.mostrar("Waiting....",snackBar.simbCargando);
    fetch("",{credentials:'same-origin',mode:'same-origin',method:"post",body:send})
    .then(resp=>{
      snackBar.menSimbTime("Changed!",snackBar.simbCheck,3);
    }).catch(err=>{
      console.log(err);
    })
    nameInput.addEventListener("input",this.showAll,{once:true});
  }
  hideAll(e){
    e.preventDefault();e.stopPropagation();
    this.cancelBtn.style.display="none";this.enviarBtn.style.display="none";
  }
  showAll(e){
    e.preventDefault();e.stopPropagation();
    this.cancelBtn.style.display="flex";this.enviarBtn.style.display="flex";
  }
}
let botonsActn=new Actions();