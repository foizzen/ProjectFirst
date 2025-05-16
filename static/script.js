let slider = document.querySelector('.posters-slider')
const nextButton = document.querySelector('.next-button')
const previousButton = document.querySelector('.previous-button')

const openPopUp = document.querySelector('#open_pop-up')
const closePopUp = document.querySelector('#close_pop-up')
const popUp = document.querySelector('#popUP')

nextButton.addEventListener('click', (event) => {
    const firstImage = slider.querySelector('a:first-child')
    slider.append(firstImage)

})

previousButton.addEventListener('click', (event) => {
    const lastImage = slider.querySelector('a:last-child')
    slider.prepend(lastImage)
})

openPopUp.addEventListener('click', function(event){
    event.preventDefault()

    popUp.classList.toggle('is-active')
})


closePopUp.addEventListener('click', (event) => {
    popUp.classList.toggle('is-active')
})
