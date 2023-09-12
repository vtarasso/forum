let s_letters = "qwertyuiopasdfghjklzxcvbnm"; // Буквы в нижнем регистре
let b_letters = "QWERTYUIOPLKJHGFDSAZXCVBNM"; // Буквы в верхнем регистре
let digits = "0123456789"; // Цифры
let specials = "!@#$%^&*()_-+=?|/.,:;{}"; // Спецсимволы
           
let password__input = document.getElementById('password__input');//получаем поле
let password__check = document.getElementById('password__check');//получаем блок с индикатором
       
password__input.addEventListener('keyup', function(evt){
let password__input_val = password__input.value;//получаем значение из поля
        
let is_s = false; // Есть ли в пароле буквы в нижнем регистре
let is_b = false; // Есть ли в пароле буквы в верхнем регистре
let is_d = false; // Есть ли в пароле цифры
let is_sp = false; // Есть ли в пароле спецсимволы
        
for (let i = 0; i < password__input_val.length; i++) {
    /* Проверяем каждый символ пароля на принадлежность к тому или иному типу */
    if (!is_s && s_letters.indexOf(password__input_val[i]) != -1) {
        is_s = true
    }
    else if (!is_b && b_letters.indexOf(password__input_val[i]) != -1) {
        is_b = true
    }
    else if (!is_d && digits.indexOf(password__input_val[i]) != -1) {
        is_d = true
    }
    else if (!is_sp && specials.indexOf(password__input_val[i]) != -1) {
        is_sp = true
    }
}

let rating = 0;
if (is_s) rating++; // Если в пароле есть символы в нижнем регистре, то увеличиваем рейтинг сложности
if (is_b) rating++; // Если в пароле есть символы в верхнем регистре, то увеличиваем рейтинг сложности
if (is_d) rating++; // Если в пароле есть цифры, то увеличиваем рейтинг сложности
if (is_sp) rating++; // Если в пароле есть спецсимволы, то увеличиваем рейтинг сложности
/* Далее идёт анализ длины пароля и полученного рейтинга, и на основании этого готовится текстовое описание сложности пароля */
if (password__input_val.length < 6 && rating < 3) {
    password__check.style.width = "10%";
    password__check.style.backgroundColor = '#e7140d';
}
else if (password__input_val.length < 6 && rating >= 3) {
    password__check.style.width = "50%";
    password__check.style.backgroundColor = '#ffdb00';
}
else if (password__input_val.length >= 21 && rating >= 4) {
    password__check.style.width = "10%";
    password__check.style.backgroundColor = '#e7140d';
}
else if (password__input_val.length >= 8 && rating < 3) {
    password__check.style.width = "50%";
    password__check.style.backgroundColor = '#ffdb00';
}
else if (password__input_val.length >= 6 && rating == 1) {
    password__check.style.width = "10%";
    password__check.style.backgroundColor = '#e7140d';
}
else if (password__input_val.length >= 6 && rating > 1 && rating < 4) {
    password__check.style.width = "50%";
    password__check.style.backgroundColor = '#ffdb00';
}
else if (password__input_val.length >= 8 && rating == 4) {
    password__check.style.width = "100%";
    password__check.style.backgroundColor = '#61ac27';
}
});