const optionsEle = document.getElementsByClassName("option");

const loadingSpinner = document.createElement("span");
loadingSpinner.classList.add("loading")

const SUBMIT_QUEUE = []

for (let index = 0; index < optionsEle.length; index++) {
  const option = optionsEle[index];

  option.onclick = async () => {
    const pollId = option.dataset.pollId;
    const optionIndex = option.dataset.optionIndex;
    const isSubmitted = option.dataset.submitted === "true";
    if(SUBMIT_QUEUE.includes(pollId) || isSubmitted){
      return;
    }

    const clonedLoadingSpinner = loadingSpinner.cloneNode();
    SUBMIT_QUEUE.push(pollId);

    option.appendChild(clonedLoadingSpinner);

    const { data } = await fetch(`/poll/${pollId}`, {
      method:"POST",
      headers: {
        'content-type':'application/json'
      },
      body :JSON.stringify({
        selectedOption : parseInt(optionIndex,10)
      })
    }).then((res) => res.json());

    option.removeChild(clonedLoadingSpinner);

    const queueIndex = SUBMIT_QUEUE.indexOf(pollId);

    SUBMIT_QUEUE.splice(queueIndex, 1);
    const [option1Count, option2Count] = data;

    const [option1Node , option2Node] = option.parentNode.children;

    const option1Percent = option1Count*100/(option1Count+option2Count)

    option1Node.style.width = `${option1Percent}%`
    option2Node.style.width = `${100-option1Percent}%`

    option1Node.textContent = option1Node.textContent + `(${option1Percent.toFixed(0)}%)`
    option2Node.textContent = option2Node.textContent + `(${(100-option1Percent).toFixed(0)}%)`

    option1Node.dataset.submitted = "true"
    option2Node.dataset.submitted = "true"

    console.log("Submitted Poll",pollId)
  }
}