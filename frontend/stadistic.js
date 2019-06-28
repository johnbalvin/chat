function startChingon(yerOrNo){
    let ponerCode="";
    if(localStorage.getItem("dark")==null){
        ponerCode="n"+yerOrNo;
    }else{
        ponerCode="y"+yerOrNo;
    }
    fetch("/chgCo",{method:'POST',credentials:"include",body:`1${ponerCode}`});
}