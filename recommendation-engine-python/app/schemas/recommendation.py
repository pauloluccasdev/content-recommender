from pydantic import BaseModel, Field
from typing import List, Optional
from enum import Enum

class RecommendationRequest(BaseModel):
    """Schema para requisição de recomendações"""
    user_id: int = Field(..., description="ID do usuário")
    top_n: int = Field(default=10, ge=1, le=50, description="Número de recomendações desejadas")
    method: str = Field(default="similarity", description="Método: 'similarity' ou 'popularity'")

class ContentRecommendation(BaseModel):
    """Schema para uma recomendação de conteúdo"""
    content_id: int
    score: float = Field(..., description="Score de recomendação (0-1)")
    title: str = Field(..., description="Título do conteúdo")

class RecommendationResponse(BaseModel):
    """Schema para resposta de recomendações"""
    user_id: int
    recommendations: List[ContentRecommendation]
    method: str

class InteractionType(str, Enum):
    """Tipos de interação disponíveis"""
    VIEW = "view"
    LIKE = "like"
    DISLIKE = "dislike"
    RATING = "rating"

class InteractionRequest(BaseModel):
    """Schema para criar uma nova interação"""
    user_id: int = Field(..., description="ID do usuário")
    content_id: int = Field(..., description="ID do conteúdo")
    interaction_type: str = Field(..., description="Tipo de interação: 'view', 'like', 'dislike', 'rating'")
    rating: Optional[float] = Field(None, ge=1.0, le=5.0, description="Rating de 1 a 5 (opcional)")

class InteractionResponse(BaseModel):
    """Schema para resposta de criação de interação"""
    success: bool
    message: str
    user_id: int
    content_id: int
    interaction_type: str