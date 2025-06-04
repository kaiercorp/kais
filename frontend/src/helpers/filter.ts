export const customFilterTrials = (filter: any, trials: any) => {
  let newTrials = trials.filter((trial: any) => {
    const accuracy = trial.accuracy * 100
    if (accuracy < filter.accuracy.min || accuracy > filter.accuracy.max) {
      return false
    }

    let inference_time = trial.inference_time * 1000
    if (inference_time < filter.inference_time.min || inference_time > filter.inference_time.max) {
      return false
    }

    if (!trial.state) {
      return false
    }

    if (filter.state[trial.state] === false) {
      return false
    }

    if (trial.created_at) {
      const trialDate = new Date(trial.created_at)

      let startDate = new Date(filter.startDate)
      if (trialDate.getFullYear() < startDate.getFullYear()) {
        return false
      } else if (
        trialDate.getFullYear() === startDate.getFullYear() &&
        trialDate.getMonth() < startDate.getMonth()
      ) {
        return false
      } else if (
        trialDate.getMonth() === startDate.getMonth() &&
        trialDate.getDate() < startDate.getDate()
      ) {
        return false
      }

      let endDate = new Date(filter.endDate)
      if (trialDate.getFullYear() > endDate.getFullYear()) {
        return false
      } else if (
        trialDate.getFullYear() === endDate.getFullYear() &&
        trialDate.getMonth() > endDate.getMonth()
      ) {
        return false
      } else if (trialDate.getMonth() === endDate.getMonth() && trialDate.getDate() > endDate.getDate()) {
        return false
      }
    }

    return true
  })

  return newTrials
}
