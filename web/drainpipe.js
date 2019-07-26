
var downspout = (() => {
    let token = "abc-def-ghi";
    let user = "geoff@gerrietts.net";
    let url = "localhost:9999";
    let loadImg = (evt) => {
        let d = {
            token: token,
            user: user,
            event: evt,
        }
        let j = JSON.stringify(d);
        let q = atob(j);
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

document.onload(() => {
    downspout.LoadImg("L")
});
document.onunload(() => {
    downspout.LoadImg("U")
});