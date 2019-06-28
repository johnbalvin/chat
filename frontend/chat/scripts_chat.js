{
    if(document.body.dataset.myid==""){
        document.querySelector("#newMessage").style.display="none";
        document.querySelector("#perfilWrapper").style.display="none";
    }else{
        document.querySelector("#acces").style.display="none";
    }
}
class Chat{
    constructor(){
        this.sendMessage=this.sendMessage.bind(this);
        this.searchMessagesBefore=this.searchMessagesBefore.bind(this);
        this.messages=document.querySelector("#messages");
        this.messagesWrapper=document.querySelector("#messagesWrapper");
        this.waiting=document.querySelector("#waiting");
        this.tmpl=document.querySelector("#tmplMesage").content.querySelector(".message");
        this.me={id:"",name:"",photo:""};
        this.usersData={};
        this.inputMessage=document.querySelector("#newMessage textarea");
        this.searchMessages=true;
        this.start();
    }
    start(){
        const btnSend=document.querySelector("#newMessage button");
        this.me.id=document.querySelector("body").dataset.myid;
        this.me.name=document.querySelector("#nameMenu1").textContent;
        this.me.photo=document.querySelector("#fotoMenu1 img").src;
        btnSend.addEventListener("click",this.sendMessage);
        this.messagesWrapper.addEventListener("scroll",this.searchMessagesBefore,{passive: true});
        const messages=document.querySelectorAll("#messages .message .userDate");
        for(let i=0,tam=messages.length;i<tam;i++){
            const message= messages[i];
            const date=new Date(parseInt(message.dataset.when)/1000000);
            message.textContent=date.toLocaleString();
        }
        this.messagesWrapper.scrollTop=this.messages.getBoundingClientRect().height;
        this.inputMessage.focus();
    }
    sendMessage(e){
        const nuevo=this.tmpl.cloneNode(true);
        nuevo.querySelector(".userText").textContent=this.inputMessage.value;
        nuevo.querySelector(".userImg").src=this.me.photo;
        nuevo.querySelector(".userName").textContent=this.me.name;
        const date=new Date();
        nuevo.querySelector(".userDate").textContent=date.toLocaleString();
        nuevo.classList.add("me");
        this.messages.appendChild(nuevo);
        this.messagesWrapper.scrollTop=this.messages.getBoundingClientRect().height; 
        const send=new FormData();
        send.append("c","n");
        send.append("m",this.inputMessage.value);
        this.inputMessage.value="";
        this.inputMessage.focus();
        fetch("",{credentials:'include',method:"POST",body:send})
        .then(resp => {
           this.waiting.style.display="none";
           switch (resp.status){
               case 200:
                    nuevo.querySelector(".spinning").remove();
                    return Promise.reject("nop")
               case 412:
                    nuevo.remove();
                    return Promise.reject("nop")    
               default:
                   return Promise.reject("server")
           }
       })
       .catch(err=>{
        if(err=="nop")return
        console.log(err);
           switch (err){
               case "server":
                   snackBar.mostrar("Server Error",snackBar.simbCross,2);
                   return
               default:
                   snackBar.mostrar("Conection error",snackBar.simbCross,2);
               return
           }
       })
    }
    searchMessagesBefore(){
        if(this.messagesWrapper.scrollTop != 0 || !this.searchMessages)return
        this.waiting.style.display="flex";
        const lastMessage=document.querySelector("#messagesWrapper .message [data-when]");
        const send=new FormData();
        if(lastMessage==null){//means no messages
            send.append("w","0");
        }else{
            send.append("w",lastMessage.dataset.when);
        }
        send.append("c","b");
        fetch("",{credentials:'include',method:"POST",body:send})
         .then(resp => {
            this.waiting.style.display="none";
            switch (resp.status){
                case 200:
                    return resp.json()
                case 204:
                    this.searchMessages=false;
                    return Promise.reject("nop")
                default:
                    return Promise.reject("server")
            }
        })
        .then(data => {
            this.putMessagesBefore(data);
        })
        .catch(err=>{
            if(err=="nop")return
            console.log(err);
            switch (err){
                case "server":
                    snackBar.mostrar("Server Error",snackBar.simbCross,2);
                    return
                default:
                    snackBar.mostrar("Conection error",snackBar.simbCross,2);
                return
            }
        })
    }
    async putMessagesBefore(data){
        this.searchMessages=data.s;
        const messages=data.m;
        let usersID ={};
        for(let i=0,tam=messages.length;i<tam;i++){
            const info=messages[i];
            if(!(info.f in this.usersData)){
                usersID[info.f]=true; //to avoid repeating users ids
            }
        }
        if (Object.keys(usersID).length>=0){
            await this.searchUsersData(usersID);
        }
        for(let i=0,tam=messages.length;i<tam;i++){
            const nuevo=this.tmpl.cloneNode(true);
            const info=messages[i];
            const infoUser=this.usersData[info.f];
            nuevo.querySelector(".userImg").src=infoUser.p;
            nuevo.querySelector(".userName").textContent=infoUser.n;
            if(this.me.id==info.f){
                nuevo.classList.add("me");
            }
            const date=new Date(parseInt(info.w)/1000000);
            nuevo.querySelector(".userText").textContent=info.t;
            nuevo.querySelector(".userDate").textContent=date.toLocaleString();
            nuevo.querySelector(".spinning").remove();
            this.messages.prepend(nuevo);
        }
        this.messages.prepend(this.waiting);
    }
   async searchUsersData(usersID){
        const send=new FormData();
        send.append("u",JSON.stringify(usersID));
 await  fetch("/UsersData",{credentials:'include',method:"POST",body:send})
         .then(resp => {
            switch (resp.status){
                case 200:
                    return resp.json()
                default:
                    return Promise.reject("server")
            }
        })
        .then(data => {
            for(let i=0,tam=data.length;i<tam;i++){
                const dataUser = data[i];
                this.usersData[dataUser.i]={n:dataUser.n,p:dataUser.p};
            }
        })
        .catch(err=>{
            console.log(err);
            switch (err){
                case "server":
                    snackBar.mostrar("Server Error",snackBar.simbCross,2);
                    return
                default:
                    snackBar.mostrar("Conection error",snackBar.simbCross,2);
                return
            }
        })
    }
}
new Chat();