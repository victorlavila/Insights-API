package com.leo.architecturetraining.data.repository

import com.leo.architecturetraining.data.datasource.CharacterDataSource
import com.leo.architecturetraining.data.model.CharacterResponse
import com.leo.architecturetraining.domain.repository.CharacterRepository

class CharacterRepositoryImpl(
    private val characterDataSource: CharacterDataSource
) : CharacterRepository {
    override suspend fun getAllCharacter(): CharacterResponse {
        return characterDataSource.getAllCharacter()
    }
}

