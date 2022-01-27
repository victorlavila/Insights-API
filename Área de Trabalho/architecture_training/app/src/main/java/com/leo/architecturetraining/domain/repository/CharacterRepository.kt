package com.leo.architecturetraining.domain.repository

import com.leo.architecturetraining.data.model.CharacterResponse

interface CharacterRepository {
    suspend fun getAllCharacter(): CharacterResponse
}