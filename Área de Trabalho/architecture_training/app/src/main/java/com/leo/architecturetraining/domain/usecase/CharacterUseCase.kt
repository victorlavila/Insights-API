package com.leo.architecturetraining.domain.usecase

import com.leo.architecturetraining.data.model.CharacterResponse
import com.leo.architecturetraining.domain.repository.CharacterRepository

class CharacterUseCase(
    private val characterRepository: CharacterRepository
) {
   suspend fun getAllCharacter(): CharacterResponse {
       return characterRepository.getAllCharacter()
   }
}