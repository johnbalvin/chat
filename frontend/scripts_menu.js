class Menu{
  constructor(){
    this.stopPropagation=this.stopPropagation.bind(this);
    this.hide=this.hide.bind(this);
    this.mostrar=this.mostrar.bind(this);
    this.optionshand=document.querySelector("#optionsNow");
    this.menuBtn=document.querySelector("#menuBtn");
    this.comenzar();
  }
  comenzar(){
    this.menuBtn.addEventListener("click",this.mostrar);
  }
  stopPropagation(e){
    e.stopPropagation();
  }
  mostrar(e){
    if(this.menuBtn.dataset.active=="true"){
      this.hide();
      return
    }
    this.menuBtn.dataset.active="true";
    this.optionshand.style.display="grid";
    this.optionshand.addEventListener("click",this.stopPropagation);
    requestAnimationFrame(()=>{
      requestAnimationFrame(()=>{
        this.optionshand.style.transform="translateX(300px)";
        document.body.addEventListener("click",this.hide,{once:true});
      });
    });
  } 
  hide(e){
    this.optionshand.addEventListener("transitionend",()=>{
      this.optionshand.style="none";
    },{once:true});
    this.optionshand.style.transform="translateX(0)";
    this.menuBtn.dataset.active="false";
    this.optionshand.removeEventListener("click",this.stopPropagation);
    document.body.removeEventListener("click",this.hide);
  }
}
const menu=new Menu();