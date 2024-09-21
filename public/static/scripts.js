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
  const form = document.querySelector('#form-fields')
  const group = document.querySelector('#group-field')
  const debut = document.querySelector('#debut-field')
  const groupType = document.querySelector('#group-type-field')

  form.style.display = 'block'
  if (type === 'group') {
    group.style.display = 'none'
    debut.style.display = 'flex'
    groupType.style.display = 'flex'
  } else {
    group.style.display = 'flex'
    debut.style.display = 'none'
    groupType.style.display = 'none'
  }
}
