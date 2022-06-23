import io
import os
from PIL import Image
from typing import Tuple

from langdetect import detect
import pytesseract
import yake

from message_types import EvidenceCreatedBody
from services import AShirtRequestsService

ExtractedKeyword = Tuple[str, float]


def process_content(body: EvidenceCreatedBody):
    ashirt_svc = AShirtRequestsService(
        os.environ.get('ASHIRT_BACKEND_URL', ''),
        os.environ.get('ASHIRT_ACCESS_KEY', ''),
        os.environ.get('ASHIRT_SECRET_KEY', '')
    )

    # gather content
    evidence_content = ashirt_svc.get_evidence_content(
        body.operation_slug, body.evidence_uuid, 'media')
    img = Image.open(io.BytesIO(evidence_content))

    # Run tesseract
    text = pytesseract.image_to_string(
        img, config='--oem 3 --psm 12 -c thresholding_method=2')

    # Extract keywords
    keywords = _extract_keywords(text)

    # convert keywords into a simple string
    result = ''
    for kw in keywords:
        result += f"{kw[0]} {kw[1]}\n"
    return result


def _extract_keywords(text: str) -> list[ExtractedKeyword]:
    language = detect(text)
    max_ngram_size = 2
    deduplication_thresold = 0.9
    deduplication_algo = 'jaro'
    window_size = 2
    num_keywords = 20

    custom_kw_extractor = yake.KeywordExtractor(
        lan=language,
        n=max_ngram_size,
        dedupLim=deduplication_thresold,
        dedupFunc=deduplication_algo,
        windowsSize=window_size,
        top=num_keywords,
        features=None
    )

    return custom_kw_extractor.extract_keywords(text)
