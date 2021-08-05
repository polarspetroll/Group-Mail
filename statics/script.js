document.getElementById('submit').addEventListener('click', () => {
    var p = document.getElementById('response');
    let inputs = document.getElementsByTagName('input');
    let emails = document.getElementById('emails').value;
    let bodytext = document.getElementById('text').value;
    let postdata = {
        from: inputs[0].value,
        subject: inputs[1].value,
        to: emails,
        body: bodytext
    };

    req = new XMLHttpRequest();
    req.open('POST', "/api/mail");
    req.setRequestHeader('content-type', 'application/json');
    req.send(JSON.stringify(postdata));

    req.onprogress = () => {
        p.style["color"] = "#e9ed72";
        p.innerText = "Sending...";
    }

    req.onload = () => {
        let sleep;
        let res = req.responseText;
        res = JSON.parse(res);
        if (res.status == 200) {
            p.style["color"] = "#58ed8c";
            sleep = 2000;
        } else {
            p.style["color"] = "red";
            sleep = 5000;
        }
        p.innerText = res["status"];
        setTimeout(() => {
            p.innerText = null;
        }, sleep);
    }

});


function changeDir(t) {
    let tag = document.getElementById('text');
    if (t.checked) {
        tag.setAttribute('dir', 'rtl');
    } else {
        tag.setAttribute('dir', 'ltr');
    }
}
