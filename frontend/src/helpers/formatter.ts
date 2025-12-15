// Functions that format various types of data for display and storing 

// Function that formats topic name so that the first letter is capitalized and the rest are lowercase
function formatTopicName(name: string) {
    if (!name) return ""

    name = name.trim().toLowerCase()

    return name.charAt(0).toUpperCase() + name.slice(1)
}

export {
    formatTopicName
}
