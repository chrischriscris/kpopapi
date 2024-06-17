function handleClickOnRandom() {
  document.getElementById('random-icon').classList.add('animate-spin')
  setTimeout(
    () =>
      document.getElementById('random-icon').classList.remove('animate-spin'),
    1000
  )
}
