function Prompt() {

    let toast = (c) => {
        const {
            title = "",
            icon = "success",
            position = "top-end",
        } = c

        const Toast = Swal.mixin({
            toast: true,
            title: title,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }

    let success = (c) => {
        const {
            title = "",
            text = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'success',
            title: title,
            text: text,
            footer: footer,
        })
    }

    let error = (c) => {
        const {
            title = "",
            text = "",
            footer = "",
        } = c

        Swal.fire({
            icon: 'error',
            title: title,
            text: text,
            footer: footer,
        })
    }

    async function custom(c) {
        const {
            icon = "",
            title = "",
            msg = "",
            showConfirmButton = true,
        } = c

        const { value: result } = await Swal.fire({
            icon: icon,
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            showConfirmButton: showConfirmButton,
            willOpen: () => {
                if (c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            didOpen: () => {
                if (c.didOpen !== undefined) {
                    c.didOpen();
                }
            }
        })

        if (result) {
            if (result.dismiss !== Swal.DismissReason.cancel) {
                if (result.value !== "") {
                    if(c.callback !== undefined) {
                        c.callback(result);
                    }
                } else {
                    c.callback(false);
                }
            } else {
                c.callback(false);
            }
        }

    }

    return {
        toast: toast,
        success: success,
        error: error,
        custom: custom,
    }
}

function BookRoom(id, csrf) {
    document.getElementById("check-availability-button").addEventListener("click", () => {
        let html = `
                <form id="check-availability-form" action="" method="post" class="needs-validation" novalidate>
                <div class="container">
                 <div class="row">
                    <div class="col">
                        <div class="row overflow-hidden" id="reservation-dates-modal">
                            <div class="col mt-1 mb-1">
                                <input disabled required type="text" class="form-control" name="start" id="start" placeholder="Arrival">
                            </div>
                            <div class="col mt-1 mb-1">
                                <input disabled required type="text" class="form-control" name="end" id="end" placeholder="Departure">
                            </div>
                        </div>
                    </div>
                 </div>
                 </div>
                </form>
            `

        attention.custom({
            msg: html,
            title: "Choose your dates",
            willOpen: () => {
                const elem = document.getElementById('reservation-dates-modal')
                const rp = new DateRangePicker(elem, {
                    format: 'yyyy-mm-dd',
                    showOnFocus: true,
                    buttonClass: 'btn',
                    minDate: new Date(),
                    autohide: true,
                    orientation: 'top',
                })
            },
            didOpen: () => {
                document.getElementById('start').removeAttribute('disabled')
                document.getElementById('end').removeAttribute('disabled')
            },
            callback: function(result) {
                console.log("called");

                let form = document.getElementById("check-availability-form");
                let formData = new FormData(form);
                formData.append("csrf_token", csrf);
                formData.append("room_id", id);

                fetch('/search-availability-json', {
                    method: "post",
                    body: formData
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.ok) {
                            attention.custom({
                                icon: 'success',
                                showConfirmButton: false,
                                msg: '<p>Room is available</p>'
                                    + '<p><a href="/book-room?id='
                                    + data.room_id
                                    + '&sd=' + data.start_date
                                    + '&ed=' + data.end_date
                                    + '" class="btn btn-primary">Book now!</a></p>'
                            })
                        } else {
                            attention.error({
                                title: "Our apologies!",
                                text: "No availability",
                            })
                        }
                    })
            }
        })
    })
}