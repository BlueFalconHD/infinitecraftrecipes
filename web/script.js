let deta;

document.addEventListener('DOMContentLoaded', () => {
    fetch('crafting_data.json')
        .then(response => response.json())
        .then(data => {
            deta = data;
            displayItems(data.items);
        });

    const modal = document.getElementById("modal");
    const span = document.getElementsByClassName("close-button")[0];

    span.onclick = function() {
        modal.style.display = "none";
    }

    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
    }
});

function displayItems(items) {
    const container = document.getElementById('items-container');
    for (const key in items) {
        const item = items[key];
        const div = document.createElement('div');
        div.className = 'item';
        div.innerHTML = `<span>${item.emoji} ${item.name}</span>`;
        div.addEventListener('click', () => showItemModal(item));
        container.appendChild(div);
    }
}

function showItemModal(item) {
    document.getElementById('item-name').textContent = item.name;
    document.getElementById('item-emoji').textContent = item.emoji;
    const recipesContainer = document.getElementById('recipes-container');
    recipesContainer.innerHTML = ''; // Clear previous recipes

    item.recipes.forEach(recipe => {
        if (recipe !== "") { // Hide any empty strings
            const badge = document.createElement('span');
            badge.className = 'recipe-badge';
            const recipeParts = recipe.split('+').map(part => `${deta.items[part].emoji} ${part}`);
            badge.innerHTML = recipeParts.join(' + '); // Adjust to include emojis
            recipesContainer.appendChild(badge);
        }
    });
    document.getElementById('modal').style.display = 'block';
}
