package com.leo.architecturetraining.data.datasource

import com.leo.architecturetraining.data.model.CharacterResponse

interface CharacterDataSource {
    suspend fun getAllCharacter() : CharacterResponse
}