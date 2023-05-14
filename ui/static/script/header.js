let submenu=document.getElementById("subMenu");
        function togglemenu(){
            submenu.classList.toggle("open-menu");
        }
        var dropdownBtns = document.querySelectorAll(".dropbtn");

            // Loop through the dropdown buttons and add a click event listener to each one
            dropdownBtns.forEach(function(btn) {
            btn.addEventListener("click", function() {
                // Toggle the visibility of the dropdown content
                this.nextElementSibling.classList.toggle("show");
            });
            });

            // Close the dropdown when the user clicks outside of it
            window.addEventListener("click", function(event) {
            dropdownBtns.forEach(function(btn) {
                if (!btn.contains(event.target)) {
                btn.nextElementSibling.classList.remove("show");
                }
            });
            });
            function showModal(id) {
                // document.querySelector('.modal').style.display = 'flex'
                document.getElementById(id).style.display = 'flex'
            }

            function closeModal(id) {
                // document.querySelector('.modal').style.display = 'none'
                document.getElementById(id).style.display = 'none'

            }
            
            const modals = document.querySelectorAll(".modal-content")

            modals.forEach(element => {
                element.addEventListener("click", (e) => { e.stopPropagation() })
            });