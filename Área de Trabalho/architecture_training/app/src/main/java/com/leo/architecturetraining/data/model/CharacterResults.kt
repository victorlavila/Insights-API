package com.leo.architecturetraining.data.model

data class CharacterResults(
    val id: Int,
    val name: String,
    val status: String,
    val specie: String,
    val type: String,
    val gender: String,
    val origin: Origin,
    val location: Location,
    val image: String,
    val episode: String,
    val url: String,
    val created: String,
)
