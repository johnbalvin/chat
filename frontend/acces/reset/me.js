class SingUp{
	constructor(){
		this.showThis=this.showThis.bind(this);
		this.sendInfo=this.sendInfo.bind(this);
		this.me=document.querySelector("#myForm");
		this.email=this.me.querySelector(".email");
		this.code=this.me.querySelector(".code");
		this.pass=this.me.querySelector(".pass");
		this.animation1=false;
		this.animating2=false;
		this.sussefull=false;
		this.start();
	}
	start(){
		const forms=document.querySelectorAll("#myForm form");
		for(let i=0,tam=forms.length;i<tam;i++){
			let me=forms[i];
			me.addEventListener("submit",e=>{
				e.preventDefault();
				let currentStep=parseInt(this.me.dataset.status);
				let current = document.querySelector(`#myForm .form[data-from="${currentStep}"]`);
				let next = current.nextElementSibling;
				this.showThis(current,next,currentStep);
			});
		}
		const inputShowOrNot=document.querySelector("#myForm .showOrNot input");
		inputShowOrNot.addEventListener("change",()=>{
			if(inputShowOrNot.checked){
				this.pass.type="text";
			}else{
				this.pass.type="password";
			}
		});
	}
	async showThis(current ,next,step){
		if(this.animating1 || this.animating2) return false;
		this.animating1 = true;
		this.animating2=true;
		await this.sendInfo(step);
		if(!this.sussefull){
			this.animating1=false;
			this.animating2=false;
			return
		}
		this.me.dataset.status=step+1;
		next.addEventListener("transitionend",()=>{this.animating1=false},{once:true});
		next.style.display="block";
		requestAnimationFrame(()=>{
			requestAnimationFrame(()=>{
				next.style.opacity=1;
				next.style.transform="scale(1) translateX(0)";
			});
		});
		current.style.position="absolute";
		next.style.position="relative";
		current.addEventListener("transitionend",()=>{
			current.style.display="none";
			current.style.transform="scale(0) translate(0)";
			this.animating2=false;
		},{once:true});
		current.style.transform="scale(0) translateX(250%)";
		current.style.opacity=0;
		document.querySelector(`#progressbar li[data-to="${next.dataset.from}"]`).classList.add("active");
		switch (step){
			case 1:
				this.code.focus();
				break;
			case 2:
				this.pass.focus();
				break;		
		}
	}
	async sendInfo(step){
		let email=this.email.value;
		let code=this.code.value;
		let pass=this.pass.value;
		let send=new FormData();
		send.append("email",email);
		send.append("s",step);
		this.sussefull=false;
		snackBar.mostrar("Waiting...",snackBar.simbCargando);
		switch (step){
			case 1:
		await	fetch("",{method:"POST",body:send,credentials:'include'})
				.then(resp =>{
					switch (resp.status){
					case 204:
						snackBar.menSimbTime("User not in database",snackBar.simbCross,5);
						break;
					case 500:
						snackBar.menSimbTime("Server error",snackBar.simbCross,5);
						break;   
					case 200:
						snackBar.menSimbTime("",snackBar.simbCargando,0);
						this.sussefull=true;
						break;  
					}
				}).catch(err=> {
					console.log(err);
					snackBar.menSimbTime("Connection error, try again",snackBar.simbCross,3)
				});
				break;
			case 2:
				send.append("code",code);
		await	fetch("",{method:"POST",body:send,credentials:'include'})
				.then(resp =>{
					switch (resp.status){
					case 409:
						snackBar.menSimbTime("Expired code".snackBar.simbCross,5);
						break;
					case 500:
						snackBar.menSimbTime("Server error",snackBar.simbCross,5);
						break;   
					case 200:
						snackBar.menSimbTime("",snackBar.simbCargando,0);
						this.sussefull=true;
						break;  
					}
				}).catch(err=> {
					console.log(err);
					snackBar.menSimbTime("Connection error, try again",snackBar.simbCross,4)
				});
				break;
			case 3:
				send.append("code",code);
				send.append("pass",pass);
		await	fetch("",{method:"POST",body:send,credentials:'include'})
				.then(resp =>{
					switch (resp.status){
					case 409:
						snackBar.menSimbTime("El código expiró".snackBar.simbCross,5);
						break;
					case 500:
						snackBar.menSimbTime("Error del servidor",snackBar.simbCross,5);
						break;   
					case 200:
						window.location.href = "/"; 
						snackBar.menSimbTime("",snackBar.simbCargando,0);
						break;
					}
				}).catch(err=> {
					console.log(err);
					snackBar.menSimbTime("Connection error, try again",snackBar.simbCross,2)
				});
				break;    
		}
	}
}
new SingUp();