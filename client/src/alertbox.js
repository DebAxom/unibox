export function createAlertBox({
    title = "Confirm Action",
    message = "Are you sure you want to continue ?",
    btn_a = "Confirm",
    btn_b = "Cancel"
} = {}) {

    return new Promise((resolve) => {
        const fragment = document.createDocumentFragment();
        const alertbox = document.getElementById("alertbox");
        alertbox.style.display = "flex";

        const box = document.createElement("div");
        box.className = "box";

        const h1 = document.createElement("h1");
        h1.textContent = title;

        const p = document.createElement("p");
        p.textContent = message;

        const textarea = document.createElement("textarea");
        textarea.maxLength = 100;
        textarea.rows = 3;
        textarea.placeholder = "Enter text (max 100 characters)";

        const btns = document.createElement("div");
        btns.className = "btns";

        const btnA = document.createElement("button");
        btnA.className = "btn-a";
        btnA.textContent = btn_a;

        const btnB = document.createElement("button");
        btnB.className = "btn-b";
        btnB.textContent = btn_b;

        function close(action) {
            const msg = textarea.value;

            alertbox.innerHTML = "";
            alertbox.style.display = "none";

            resolve({ action, msg });
        }

        btnA.addEventListener("click", () => close(btn_a));
        btnB.addEventListener("click", () => close(btn_b));

        btns.appendChild(btnA);
        btns.appendChild(btnB);

        box.appendChild(h1);
        box.appendChild(p);
        box.appendChild(textarea);
        box.appendChild(btns);

        fragment.appendChild(box);
        alertbox.appendChild(fragment);
    });
}