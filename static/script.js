let slider = document.querySelector('.posters-slider')
const nextButton = document.querySelector('.next-button')
const previousButton = document.querySelector('.previous-button')


const popUpLogin = document.querySelector('#popUP')
const openPopUpLogin = document.querySelector('#open_pop-up')
const closePopUpLogin = document.querySelector('#close_pop-up')
const loginButton = document.querySelector('#login-button')
const loginFormElement = document.querySelector('#pop_up_form-login')

const popUpRegistration = document.querySelector('#popUP_reg')
const openPopUpRegistration = document.querySelector('#open_reg_pop-up')
const closePopUpRegistration = document.querySelector('#close_pop-up-reg')
const registrationButton = document.querySelector('#registration-button')
const registrationFormElement = document.querySelector('#pop_up_form-registration')


if (localStorage.username === 'is-logged-in') {
    const userNickname = 
        (typeof formDataObjectFetch !== 'undefined' && formDataObjectFetch.username) || 
        (typeof formDataObject !== 'undefined' && formDataObject.username) || ""


    openPopUpLogin.classList.add('has-username');
    openPopUpLogin.setAttribute('data-title', userNickname);
}




nextButton.addEventListener('click', (event) => {
    const firstImage = slider.querySelector('a:first-child')
    slider.append(firstImage)

})

previousButton.addEventListener('click', (event) => {
    const lastImage = slider.querySelector('a:last-child')
    slider.prepend(lastImage)
})



openPopUpLogin.addEventListener('click', function(event){
    event.preventDefault()

    popUpLogin.classList.add('is-active')
})

closePopUpLogin.addEventListener('click', (event) => {
    popUpLogin.classList.remove('is-active')
})


openPopUpRegistration.addEventListener('click', function(event){
    event.preventDefault()

    popUpRegistration.classList.add('is-active')
    popUpLogin.classList.remove('is-active')
})

closePopUpRegistration.addEventListener('click', (event) => {
    popUpRegistration.classList.remove('is-active')
})

// -- Валидация 
class FormsValidation {
    selectors = {
        form: '[data-js-form]',
        fieldErrors: '[data-js-form-field-errors]',
    }

    errorMessages = {
        valueMissing: () => 'Please fill this field',
        patternMismatch:  ({ title }) => title || 'Data type is incorrect',
        tooShort: ({ minLength }) => `Too short. Minimum ${minLength} characters required.`,
        tooLong: ({ maxLength }) => `Too long. Maximum ${maxLength} characters required.`,
    }

    constructor () {
        this.bindEvents()
    }

    manageErrors (fieldControlElement, errorMessages) {
        const fieldErrorsElement = fieldControlElement.parentElement.querySelector(this.selectors.fieldErrors)

        fieldErrorsElement.innerHTML = errorMessages
        .map( (message) => `<span class = "field__error">${message}</span>`)
        .join('')
    }

    validateField(fieldControlElement) {
        const errors = fieldControlElement.validity
        const errorMessages = []

        Object.entries(this.errorMessages).forEach(([errorType, getErrorMessage]) =>{
            if(errors[errorType]) {
                errorMessages.push(getErrorMessage(fieldControlElement))
            }
        })

        this.manageErrors(fieldControlElement, errorMessages)

        const isValid = errorMessages.length === 0

        fieldControlElement.ariaInvalid = !isValid

        return isValid
    }

    onBlur(event) {
        const { target } = event
        const isFormField = target.closest(this.selectors.form)
        if(isFormField) {
            this.validateField(target)
        }
    }

    onSubmit(event) {
        const isFormElement = event.target.matches(this.selectors.form)

        if(!isFormElement){
            return
        }

        const requiredControlElements = [...event.target.elements].filter((element)=> element.required)
        let isFormValid = true  
        let firstInvalidFieldControl = null


        requiredControlElements.forEach((element)=> {
            const isFieldValid = this.validateField(element)

            if(!isFieldValid) {
                isFormValid = false

                if(!firstInvalidFieldControl){
                    firstInvalidFieldControl = element
                }
            }
        })

        if (!isFormValid) {
            event.preventDefault()
            firstInvalidFieldControl.focus()
        }
    }

    bindEvents() {
        document.addEventListener('blur', (event) => {
            this.onBlur(event)
        }, {capture: true})
        document.addEventListener('submit', (event) => this.onSubmit(event))
    }
}

new FormsValidation()


loginFormElement.addEventListener('submit', (event) => {
    event.preventDefault()
    const formData = new FormData(loginFormElement)
    const formDataObject = Object.fromEntries(formData)
    const buttonError = document.querySelector('#login-password-error')

    fetch('/login', {
        method: 'post',
        body: JSON.stringify(formDataObject),
        headers: {
            'Content-Type': 'application/json'
        },
    }).then((response) => {

        if(response.status != 200){
            console.log(response.text)
            buttonError.innerHTML = `<h3 class="pattern-mismatch">Incorrect username or password</h3>`

            return
        }
        

        localStorage.setItem('username', 'is-logged-in')

        return response.json()
    }).catch((error) =>{ 
        console.log('Error:', error)
    })
    return formDataObject.username
})


registrationFormElement.addEventListener('submit', (event) => {
    event.preventDefault()

    const formData = new FormData(registrationFormElement)
    const formDataObject = Object.fromEntries(formData)
    const buttonError = document.querySelector('#registration-password-error')

    if(formDataObject.password1 != formDataObject.password2){
        console.log('Something wrong')
        
        buttonError.innerHTML = `<h3 class="pattern-mismatch">Passwords don't match</h3>`
        return
    }

    const loginPassword = formDataObject.password1
        
    const formDataObjectFetch = {
        username: `${formDataObject.username}`,
        password:`${loginPassword}`,
    }

    fetch('/reg', {
        method: 'post',
        body: JSON.stringify(formDataObjectFetch),
        headers: {
            'Content-Type': 'application/json'
        },
    }).then((response) => {
        console.log('response:', response)

        localStorage.setItem('username', 'is-logged-in')

        return response.json()
    }).then ((json) => {

        console.log('json:', json)
    }).catch((error) =>{ 
        console.log('Error:', error)
    })

    return formDataObjectFetch.username
})


