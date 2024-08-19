function handleClickOnRandom() {
  document.getElementById('random-icon').classList.add('animate-spin')
  setTimeout(
    () =>
      document.getElementById('random-icon').classList.remove('animate-spin'),
    1000
  )
}

/**
 * Handle the change of the type of addition
 * @param {HTMLInputElement} type
 */
function handleTypeChange(type) {
  const group = document.querySelector('#group-field')
  const debut = document.querySelector('#debut-field')

  if (type === 'group') {
    group.style.display = 'none'
    debut.style.display = 'block'
  } else {
    group.style.display = 'block'
    debut.style.display = 'none'
  }
}
