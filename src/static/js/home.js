requestAnimationFrame(() => {
    setCardsWidth();
});

function setCardsWidth() {
    const recContainer = document.getElementById("course-cards-row-rec");
    const newContainer = document.getElementById("course-cards-row-new");
    const recCards = recContainer.querySelectorAll(".card");
    const newCards = newContainer.querySelectorAll(".card");

    // Set the width of the cards to be the same
    const recCardWidth = calculateCardWidth(recContainer);
    setCardWidth(recCards, recCardWidth);

    const newCardWidth = calculateCardWidth(newContainer);
    setCardWidth(newCards, newCardWidth);
}

function calculateCardWidth(container) {
    const containerWidth = container.clientWidth - 8; // Subtracting 8px for the padding of the container
    const minCardWidth = 200; // Minimum width for each card
    const gap = 4; // Gap between cards

    // Calculate the width of each card based on the container width and number of visible cards
    const visibleCards = Math.max(Math.floor((containerWidth + gap) / (minCardWidth + gap)), 1);
    return (containerWidth - (visibleCards - 1) * gap) / visibleCards;
}

function calculateVisibleCardsRec() {
    const containerWidth = document.getElementById("course-cards-row-rec").clientWidth - 8; // Subtracting 8px for the padding
    const minCardWidth = 200; // Minimum width for each card
    const gap = 4; // Gap between cards

    return Math.max(Math.floor((containerWidth + gap) / (minCardWidth + gap)), 1);
}

function setCardWidth(cards, width) {
    cards.forEach(card => {
        card.style.width = `${width}px`;
        card.style.minWidth = `${width}px`;
    });
}