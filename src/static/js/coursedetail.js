// uncheck checkbox by id
function uncheck(id) {
    const checkbox = document.getElementById(id);
    console.log('uncheck', id, checkbox);
    if (checkbox) {
        const clickEvent = new MouseEvent('click', {
            bubbles: true,
            cancelable: true,
        });
        checkbox.dispatchEvent(clickEvent);
    }

    const surveyFilters = document.getElementById('survey-filters-form');
    if (surveyFilters) {
        const updateFilters = new CustomEvent('filters-changed', {});
        surveyFilters.dispatchEvent(updateFilters);
    }
}