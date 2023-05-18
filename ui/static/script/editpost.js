const modal=document.querySelector(".modalPost")
const overlay=document.querySelector(".overlay")
const openModalBtn=document.querySelector(".btn-edit")
const closeModalBtn=document.querySelector(".btn-close")

const openModal1=function(){
    modal.classList.remove("hidden");
    overlay.classList.remove("hidden");
}
openModalBtn.addEventListener("click",openModal1);

const closeModal1=function(){
    modal.classList.add("hidden");
    overlay.classList.add("hidden");
}
closeModalBtn.addEventListener("click",closeModal1);
overlay.addEventListener("click",closeModal1);
