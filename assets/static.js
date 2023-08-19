const optionsEle = document.getElementsByClassName("option");

const POLLS_CONTAINER = document.getElementById("polls-container");
const POLL_QUESTION_NODES = Array.from(document.getElementsByClassName("poll-container"))

let fetching = false;

const loadingSpinner = document.createElement("span");
loadingSpinner.classList.add("loading")

const onOptionClick = async (option) => {
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

const SUBMIT_QUEUE = []

for (let index = 0; index < optionsEle.length; index++) {
  const option = optionsEle[index];

  option.onclick = () => onOptionClick(option)
}

const createPollQuestionNode = (poll, index) => {
  const pollContainer = document.createElement("div");
  pollContainer.classList.add("poll-container")
  pollContainer.dataset.pollId=poll.id
  
  const pollQuestion = document.createElement("div")
  pollQuestion.classList.add("poll-question");
  pollQuestion.textContent = poll.question

  const optionsContainer = document.createElement("div");
  optionsContainer.classList.add("poll-options-container");

  const option1Node = document.createElement("div")
  option1Node.textContent = poll.options[0]
  option1Node.classList.add("option","option-0")
  option1Node.dataset.pollId = poll.id;
  option1Node.data.optionIndex=0
  option1Node.onclick = () => onOptionClick(option1Node);

  const option2Node = document.createElement("div")
  option2Node.textContent = poll.options[1]
  option2Node.classList.add("option","option-1")
  option2Node.dataset.pollId = poll.id;
  option2Node.data.optionIndex=1
  option2Node.onclick = () => onOptionClick(option2Node)

  optionsContainer.appendChild(option1Node)
  optionsContainer.appendChild(option2Node)

  pollContainer.appendChild(pollQuestion)
  pollContainer.appendChild(optionsContainer);

  POLLS_CONTAINER.append(pollContainer);

  POLL_QUESTION_NODES[index] = pollContainer;
}

setInterval(() => {
 (async () => {
    if(SUBMIT_QUEUE.length > 0 || fetching){
      return;
    }
    fetching = true;
    let updatedPolls;
    try {
      updatedPolls = await fetch("/polls").then(res => res.json());
    } catch (error) {
      window.location.reload()// reload 
    }

    updatedPolls.forEach((poll, index) => {
      const pollNode = POLL_QUESTION_NODES[index];
      // Don't update if not submitted and node is already there
      if(!poll.submitted && pollNode){
        return;
      }

      if(!pollNode){
        createPollQuestionNode(poll, index)
      }

      // updaete node
      const [, optionsContainer]=pollNode.children;

      const [option1Node, option2Node] = optionsContainer.children;

      const [option1Count, option2Count] = poll.submissions;

      const option1Percent = option1Count*100/(option1Count+option2Count)

      option1Node.style.width = `${option1Percent}%`
      option2Node.style.width = `${100-option1Percent}%`

      option1Node.textContent = poll.options[0] + `(${option1Percent.toFixed(0)}%)`
      option2Node.textContent = poll.options[1] + `(${(100-option1Percent).toFixed(0)}%)`

      option1Node.dataset.submitted = "true"
      option2Node.dataset.submitted = "true"

    })

    fetching=false;

 })()
},750 + Math.floor(Math.random()*50)) // add a random value so all the client send api request at diffrerent times