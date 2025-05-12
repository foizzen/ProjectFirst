let slider = document.querySelector('.posters-slider')
const nextButton = document.querySelector('.next-button')
const previousButton = document.querySelector('.previous-button')



nextButton.addEventListener('click', (event) => {
    const firstImage = slider.querySelector('a:first-child')
    slider.append(firstImage)

})

previousButton.addEventListener('click', (event) => {
    const lastImage = slider.querySelector('a:last-child')
    slider.prepend(lastImage)
})
