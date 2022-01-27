package com.leo.architecturetraining.domain.model

import com.leo.architecturetraining.data.model.Location
import com.leo.architecturetraining.data.model.Origin

data class CharactersInfo(
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

data class ToInfo(
    val count: Int,
    val pages: Int,
    val next : String,
    val prev: String,
)