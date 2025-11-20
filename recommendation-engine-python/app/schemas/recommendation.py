from email.policy import default
from pydantic import BaseModel, Field
from typing import List

class RecommendationRequest(BaseModel):
    """Schema para requisição de recomendações"""
    user_id: str = Field(..., description="ID do usuário")
    top_n: int = Field(default=10, ge=1, le=50, description="Número de recomendações desejadas")
    method: str = Field(default="similarity", description="Método: 'similarity' ou 'popularity'")

class ContentRecommendation(BaseModel):
    """Schema para uma recomendação de conteúdo"""
    content_id: int
    score: float = Field(..., description="Score de recomendação (0-1)")
    title: str = Field(..., description="Título do conteúdo")

class RecommendationResponse(BaseModel):
    """Schema para resposta de recomendações"""
    user_id: str
    recommendations: List[ContentRecommendation]
    method: str