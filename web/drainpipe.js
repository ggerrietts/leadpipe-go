
var downspout = (() => {
    let token = "abc-def-ghi";
    let user = "geoff@gerrietts.net";
    let url = "http://localhost:8080";
    let loadImg = (evt) => {
        let d = {
            token: token,
            user: user,
            event: evt,
        }
        let j = JSON.stringify(d);
        let q = btoa(j);
        let u = url + "/img.png?q=" + q;
        let i = document.createElement('img');
        i.src = u;
        document.body.insertBefore(i, document.body.firstChild);
    }
    return {
        SetUser: (u) => { user = u; loadImg("N") },
        LoadImg: () => { loadImg() }
    }
})();

window.onload = () => {
    downspout.LoadImg("L")
};
window.onunload = () => {
    downspout.LoadImg("U")
};