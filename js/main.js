function postContactToGoogle() {
    var email = $('#user_email').val();
    var message = $('#user_message').val();
    console.log(email, message);
    var request = $.ajax({
	url: "https://docs.google.com/a/codeignition.co/forms/d/1wfr8OXaNiCSYR6E-P7NALnLjx-cS2gyNvekL82ArBv8/formResponse",
        data: {
	    "entry.596447052" : email,
	    "entry.1743652851" : message
	},
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
