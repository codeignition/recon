function postContactToGoogle() {
    var email = $('#user_email').val();
    console.log(email);
    var request = $.ajax({
	url: "https://docs.google.com/a/codeignition.co/forms/d/1tt-eSxJ1oG_3W2gU2tvdeK9SGsU4_n44DH3tNeYiZqM/formResponse",
        data: {"entry.440186902" : email},
        type: "POST",
        dataType: "xml",
	statusCode: {
                    0: showSuccess,
                    200: showSuccess,
	            404: showFailure
        }
    });
}

function showSuccess() {
    $('#landing-input').hide("slow");
    $('#landing-success-message').show("slow");
}

function showFailure() {
    $('#landing-failure-message').show("slow");
}
