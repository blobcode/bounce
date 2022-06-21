let url = "";

function submit() {
  let input = document.getElementById("input");
  let result = document.getElementById("result");
  let copy = document.getElementById("copy");

  let data = {
    url: input.value,
  };

  fetch("/new/", {
    method: "POST", // *GET, POST, PUT, DELETE, etc.
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(data), // body data type must match "Content-Type" header
  })
    .then((response) => response.json())
    .then((data) => {
      console.log(data);
      let string = `${window.location.href}r/${data.id}`;
      result.innerText = string;
      result.href = string;
      url = string;
      copy.style.display = "inline";
    });
}
function copy() {
  navigator.clipboard.writeText(url);
}
